from airflow import DAG
from airflow.operators.python import PythonOperator
from datetime import datetime


def placeholder_task():
    # TODO: Replace with ingestion + dbt tasks
    return "ok"


with DAG(
    dag_id="revenue_dashboard_placeholder",
    start_date=datetime(2024, 1, 1),
    schedule_interval="@daily",
    catchup=False,
) as dag:
    run = PythonOperator(
        task_id="placeholder",
        python_callable=placeholder_task,
    )
