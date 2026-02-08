select
    user_id,
    account_id,
    email_hash,
    country_code,
    user_type,
    cast(signup_ts as timestamp) as signup_ts,
    cast(last_active_ts as timestamp) as last_active_ts
from {{ source('app', 'users') }}
