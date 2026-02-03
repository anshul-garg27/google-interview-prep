{{ config(materialized = 'table', tags=["hourly"], order_by='profile_id') }}
with
ig_cats as (
    select
        profile_id,
        ig_category,
        tupleElement(cat_from_ig_cat[1], 1) ig_cat,
        tupleElement(cat_from_ig_cat[1], 2) ig_score
    from ds.categorization_zsl
    where ig_score > 0.9
),
bio_cats as (
    select
        profile_id,
        ig_category,
        tupleElement(cat_from_bio[1], 1) bio_cat,
        tupleElement(cat_from_bio[1], 2) bio_score
    from ds.categorization_zsl
    where bio_score  > 0.9
),
predicted_cats as (
    select
        ifNull(coalesce(if(ig_cats.profile_id = '', NULL, ig_cats.profile_id), bio_cats.profile_id), 'MISSING') profile_id,
        coalesce(if(ig_cat = '', NULL, ig_cat), bio_cat) category
    from
    ig_cats FULL OUTER JOIN  bio_cats on ig_cats.profile_id = bio_cats.profile_id
)
select * from predicted_cats