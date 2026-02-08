select
    account_id,
    account_name,
    industry,
    plan_type,
    sales_region,
    account_status,
    cast(created_ts as timestamp) as created_ts,
    cast(updated_ts as timestamp) as updated_ts
from {{ source('app', 'accounts') }}
