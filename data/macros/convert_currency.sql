{% macro convert_currency(amount, from_currency, to_currency='USD', rate_table='dim_exchange_rates') %}
    (
        {{ amount }} * (
            select rate
            from {{ ref(rate_table) }}
            where from_currency = {{ from_currency }}
              and to_currency = {{ to_currency }}
              and rate_date = current_date
        )
    )
{% endmacro %}
