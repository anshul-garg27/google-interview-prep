from datetime import datetime, timedelta


def fetch_scraping_windows(start_date, end_date, frequency:int = None):
    start = datetime.strptime(start_date, '%Y-%m-%d')
    end = datetime.strptime(end_date, '%Y-%m-%d')
    breaked_dates = []
    
    if frequency:
        current_date = start
        while current_date <= end:
            interval_start = current_date.strftime('%Y-%m-%dT00:00:00Z')
            interval_end = (current_date + timedelta(days=frequency-1)).strftime('%Y-%m-%dT23:59:59Z')
            if current_date + timedelta(days=frequency-1) > end:
                interval_end = end.strftime('%Y-%m-%dT23:59:59Z')
            breaked_dates.append((interval_start, interval_end))
            current_date += timedelta(days=frequency)
    else:
        if (end - start).days < 60:
            # Break into week-level intervals
            current_date = start
            while current_date <= end:
                week_start = current_date.strftime('%Y-%m-%dT00:00:00Z')
                week_end = (current_date + timedelta(days=6)).strftime('%Y-%m-%dT23:59:59Z')
                if current_date + timedelta(days=6) > end:
                    week_end = end.strftime('%Y-%m-%dT23:59:59Z')
                breaked_dates.append((week_start, week_end))
                current_date += timedelta(days=7)
        else:
            # Break into month-level intervals with 30 days window
            current_date = start
            while current_date <= end:
                month_start = current_date.strftime('%Y-%m-%dT00:00:00Z')
                month_end = (current_date + timedelta(days=29)).strftime('%Y-%m-%dT23:59:59Z')
                if current_date.month == end.month:
                    month_end = end.strftime('%Y-%m-%dT23:59:59Z')
                elif current_date + timedelta(days=29) > end:
                    month_end = (end + timedelta(days=1)).strftime('%Y-%m-%dT23:59:59Z')
                breaked_dates.append((month_start, month_end))
                current_date += timedelta(days=30)
    return breaked_dates
