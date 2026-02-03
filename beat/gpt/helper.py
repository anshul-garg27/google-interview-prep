from collections import defaultdict

import loguru
import numpy as np
import random
import pandas as pd


def is_data_consumable(data: dict, data_type: str):
    keys_to_check = ["18-24 male", "18-24 female", "25-34 male", "25-34 female",
                     "35-44 male", "35-44 female", "45-54 male", "45-54 female",
                     "55-64 male", "55-64 female", "65+ male", "65+ female"]

    if data_type == "base_gender":
        required_keys = ["gender"]
        for key in required_keys:
            if key not in data or not data[key]:
                return False
        if any(value.lower() in ["", "unknown"] for value in [data['gender']]):
            return False
    elif data_type == "base_localtion":
        required_keys = ["city", "state", "country"]
        for key in required_keys:
            if key not in data or not data[key]:
                return False
        if any(value.lower() in ["", "unknown"] for value in [data['city'], data['state'], data['country']]):
            return False
    elif data_type == "base_categ_lang_topics":
        required_keys = ["categories", "languages", "topics"]
        for key in required_keys:
            if key not in data or not data[key]:
                return False
    elif data_type == "audience_age_gender":
        if "audience" not in data or "age_gender" not in data['audience']:
            return False
        age_gender = data['audience']['age_gender']
        if not all(key in age_gender for key in keys_to_check):
            return False
        # Sanity check for valid percentage ranges
        if any(value < 0 or value > 1.0 for value in age_gender.values()):
            return False
        # Sanity check for overlapping age groups
        if len(set(age_gender.keys())) != len(age_gender.keys()):
            return False

        # Calculate the total percentage of males and females
        total_male_percentage = sum(age_gender[key] for key in age_gender if "male" in key)
        total_female_percentage = sum(age_gender[key] for key in age_gender if "female" in key)

        # Check if male or female ratio is less than 0.25%
        if total_male_percentage < 0.25 or total_female_percentage < 0.25:
            return False

    elif data_type == "audience_cities":
        if "audience" not in data or "cities" not in data['audience']:
            return False
        cities = data['audience']['cities']
        if not cities or len(cities) <= 5:
            return False
    elif data_type == "base_localtion_gender_lang":
        required_keys = ["gender", "city", "state", "country", "languages"]
        for key in required_keys:
            if key not in data or not data[key]:
                return False
            if any((value.lower() if isinstance(value, str) else value ) in ["", "unknown", None] for value in [data['gender'], data['city'], data['state'],
                                                                  data['country']]):
                return False
    return True


def ssd(a, b):
    a = np.array(a)
    b = np.array(b)
    dif = a.ravel() - b.ravel()
    return np.dot(dif, dif)


def gradient_descent(a, b, learning_rate=0.01, epochs=1000):
    if len(a) != len(b):
        raise ValueError("Arrays must have the same length")

    a = np.array(a)
    b = np.array(b)

    for epoch in range(epochs):
        # Compute the gradient of the loss with respect to array 'a'
        gradient = -2 * (b - a)

        # Update 'a' using the gradient
        a -= learning_rate * gradient
        # print(f"Iteration {epoch + 1} :- {a.tolist()}")

    return a.tolist()


def normalize_audience_age_gender(audience_age_gender_data, category):
    age_gender = audience_age_gender_data['audience']['age_gender']
    total_percentage = sum(age_gender.values())
    # normalization
    if total_percentage != 1.0:
        normalization_factor = 1 / total_percentage
        for age, value in age_gender.items():
            audience_age_gender_data['audience']['age_gender'][age] = value * normalization_factor

    file_path = 'gpt/age_gender_private_data.csv'
    df = pd.read_csv(file_path)
    if not category or not df['categories'].to_numpy().__contains__(category):
        category = 'Missing'
    filtered_rows = df[df['categories'] == category]
    filtered_rows_list = filtered_rows.values.tolist()
    filtered_rows_list = filtered_rows_list[0][1:]
    a = list(age_gender.values())
    b = filtered_rows_list
    result = gradient_descent(a, b, epochs=random.randint(50, 100))

    i = 0
    for age, value in age_gender.items():
        audience_age_gender_data['audience']['age_gender'][age] = round(result[i], 3)
        i += 1


def generate_random_values(ul, ll, size):
    random_values = np.random.uniform(high=ul, low=ll, size=size)
    rounded_values = np.round(random_values, 5)
    sorted_values = np.sort(rounded_values)[::-1]
    return sorted_values


def normalize_audience_cities(data):
    cities = data['audience']['cities']
    values_to_keys = defaultdict(list)
    #   group keys by their values
    for key, value in cities.items():
        values_to_keys[value].append(key)

    keys = list(values_to_keys.keys())
    # Add noise to keys with the same values
    for i, keys_list in enumerate(values_to_keys.values()):
        if len(keys_list) > 1:
            prev_value = values_to_keys[keys[i-1]][-1] if i-1 >= 0 else 0.015
            prev_value = cities[prev_value] if isinstance(prev_value, str) else prev_value
            next_value = values_to_keys[keys[i+1]][0] if i+1 < len(keys) else -0.015
            next_value = cities[next_value] if isinstance(next_value, str) else next_value
            upper_limit = min(0.015, prev_value - keys[i])
            lower_limit = max(-0.015, (keys[i] - next_value) * -1)

            if i == len(values_to_keys)-1:
                lower_limit = 0 if keys[i] - 0.015 < 0 else 0.015
            random_values = generate_random_values(upper_limit, lower_limit, len(keys_list))

            for i, key in enumerate(keys_list):
                noise = random_values[i]
                cities[key] = round(cities[key] + noise, 5)

    total_percentage = sum(cities.values())
    loguru.logger.debug(total_percentage)
    # normalize cities data
    if total_percentage != 1.0:
        normalization_factor = 1 / total_percentage
        loguru.logger.debug(normalization_factor)
        for city, value in cities.items():
            data['audience']['cities'][city] = round(value * normalization_factor, 5)
        loguru.logger.debug(data)
    normalization_factor = np.random.uniform(.65, .75)
    total_sum = 0
    for city, value in cities.items():
        normalized_value = round(value * normalization_factor, 5)
        data['audience']['cities'][city] = normalized_value
        total_sum += normalized_value
    data['audience']['cities']['others'] = round(1 - total_sum, 5)
    loguru.logger.debug(data)
