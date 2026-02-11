from __future__ import annotations

import os
import sqlite3
from datetime import datetime

from airflow import DAG
from airflow.operators.python import PythonOperator

DB_PATH = os.getenv("WAREHOUSE_DB_PATH", "/opt/airflow/api/dev.db")


def check_data_quality():
    if not os.path.exists(DB_PATH):
        print(f"[data_quality] SQLite DB not found at {DB_PATH}. Skipping checks.")
        return

    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()

    def fetch_one(query: str, params: tuple | None = None):
        cursor.execute(query, params or ())
        return cursor.fetchone()[0]

    metrics = {
        "orders_count": fetch_one("select count(*) from fact_orders"),
        "sessions_count": fetch_one("select count(*) from fact_sessions"),
        "active_users_count": fetch_one("select count(*) from fact_active_users"),
        "subscriptions_count": fetch_one("select count(*) from fact_subscriptions"),
    }

    anomalies = []
    if metrics["orders_count"] == 0:
        anomalies.append("orders_count is 0")
    if metrics["sessions_count"] == 0:
        anomalies.append("sessions_count is 0")
    if metrics["active_users_count"] == 0:
        anomalies.append("active_users_count is 0")

    print("[data_quality] Metrics:", metrics)
    if anomalies:
        print("[data_quality] Anomalies:", anomalies)
    else:
        print("[data_quality] All checks passed")

    conn.close()


def build_dag():
    with DAG(
        dag_id="revenue_dashboard_data_quality",
        start_date=datetime(2024, 1, 1),
        schedule_interval="@daily",
        catchup=False,
        tags=["data_quality"],
    ) as dag:
        PythonOperator(
            task_id="data_quality_checks",
            python_callable=check_data_quality,
        )

    return dag


dag = build_dag()
