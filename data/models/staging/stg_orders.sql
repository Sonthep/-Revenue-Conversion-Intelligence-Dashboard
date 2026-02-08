with source as (
    select *
    from {{ source('app', 'orders') }}
)

select
    order_id,
    account_id,
    user_id,
    product_id,
    cast(order_ts as timestamp) as order_ts,
    cast(date(order_ts) as date) as order_date,
    cast(gross_amount as numeric) as gross_amount,
    cast(discount_amount as numeric) as discount_amount,
    cast(refund_amount as numeric) as refund_amount,
    cast(tax_amount as numeric) as tax_amount,
    currency,
    payment_status,
    ingest_ts
from source
