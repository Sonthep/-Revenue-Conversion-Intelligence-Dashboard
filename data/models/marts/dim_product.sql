select
    product_id,
    sku,
    product_name,
    product_category,
    product_type,
    currency,
    cast(base_price as numeric) as base_price,
    is_active,
    cast(launch_date as date) as launch_date
from {{ source('app', 'products') }}
