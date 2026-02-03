import psycopg2
from datetime import datetime,timedelta
import time
from dateutil.relativedelta import relativedelta


# Define connection parameters for PostgreSQL

# pg_host = 'localhost'
# pg_port = 5432
# pg_database = 'coffee'
# pg_user = 'root'
# pg_password = 'Glamm@123'

pg_host = '172.31.2.21'
pg_port = 5432
pg_database = 'coffee'
pg_user = 'gccuser'
pg_password = 'dbBeat123UseRpr0d'

# Create a connection object for PostgreSQL
postgres_conn = psycopg2.connect(
    host=pg_host,
    port=pg_port,
    database=pg_database,
    user=pg_user,
    password=pg_password
)

postgres_cursor = postgres_conn.cursor()

def main():
    current_time = datetime.fromtimestamp(time.time())
    check_query = 'SELECT * FROM partner_usage WHERE next_reset_on < %s'
    postgres_cursor.execute(check_query, (current_time,))
    rows = postgres_cursor.fetchall()
    column_names = [desc[0] for desc in postgres_cursor.description]

    for item in rows:
        row_dict = dict(zip(column_names, item))
        if row_dict['end_date'] > datetime.fromtimestamp(time.time()):
            if row_dict['frequency'] == "daily":
                row_dict['next_reset_on'] = row_dict['next_reset_on'] + datetime.timedelta(days=1)
            elif row_dict['frequency'] == "weekly":
                row_dict['next_reset_on'] = row_dict['next_reset_on'] + datetime.timedelta(days=7)
            elif row_dict['frequency'] == "monthly":
                row_dict['next_reset_on'] = row_dict['next_reset_on'] + relativedelta(months=1)
            elif row_dict['frequency'] == "quarterly":
                row_dict['next_reset_on'] = row_dict['next_reset_on'] + relativedelta(months=3)
            elif  row_dict['frequency'] == "yearly":
                row_dict['next_reset_on'] = row_dict['next_reset_on'] + relativedelta(years=1)
            
            updated_at = datetime.fromtimestamp(time.time())
            update_query = 'UPDATE partner_usage SET next_reset_on = %s, updated_at = %s WHERE id = %s'
            postgres_cursor.execute(update_query, (row_dict['next_reset_on'], updated_at, row_dict['id']))
            postgres_conn.commit()
        else:
            delete_partner_query = 'DELETE FROM partner_usage WHERE partner_id = %s'
            postgres_cursor.execute(delete_partner_query, (row_dict['partner_id'],))
            postgres_conn.commit()

        delete_query = 'DELETE FROM partner_profile_page_track WHERE partner_id = %s'
        postgres_cursor.execute(delete_query, (row_dict['partner_id'],))
        postgres_conn.commit()

if __name__ == "__main__":
    main()