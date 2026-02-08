with orders as (
    select *
    from {{ ref('stg_orders') }}
)

select
    order_id,
    order_ts,
    order_date,
    account_id,
    user_id,
    product_id,
    gross_amount,
    coalesce(discount_amount, 0) as discount_amount,
    coalesce(refund_amount, 0) as refund_amount,
    coalesce(tax_amount, 0) as tax_amount,
    (gross_amount - coalesce(discount_amount, 0) - coalesce(refund_amount, 0) - coalesce(tax_amount, 0)) as net_amount,
    currency,
    payment_status,
    ingest_ts,
    current_timestamp as dbt_updated_at
from orders
where payment_status in ('completed', 'refunded', 'partial_refund')
