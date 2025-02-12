import pytest
from httpx import AsyncClient
from datetime import timedelta
from tests.models.task_models import TaskAnalytics

@pytest.mark.asyncio
class TestTaskAnalytics:
    async def test_time_based_analytics(self, api_client: AsyncClient, frozen_time):
        # Create task with 50% progress
        start = frozen_time - timedelta(days=2)
        due = frozen_time + timedelta(days=2)
        
        task_res = await api_client.post("/tasks", json={
            "name": "Analytics Task",
            "start_date": start.isoformat(),
            "due_date": due.isoformat(),
            "priority": 2
        })
        task_id = task_res.json()["id"]
        
        # Check analytics
        analytics_res = await api_client.get(f"/tasks/{task_id}/analytics")
        assert analytics_res.status_code == 200
        analytics = TaskAnalytics(**analytics_res.json())
        
        assert analytics.percentage_complete == pytest.approx(50.0, abs=0.1)
        assert analytics.completion_health == "okay"

    async def test_completed_task_analytics(self, api_client: AsyncClient, frozen_time):
        # Create and complete task
        task_res = await api_client.post("/tasks", json={
            "name": "Quick Task",
            "start_date": frozen_time.isoformat(),
            "due_date": (frozen_time + timedelta(hours=1)).isoformat(),
            "priority": 1
        })
        task_id = task_res.json()["id"]
        await api_client.put(f"/tasks/{task_id}/complete")
        
        # Verify analytics
        analytics_res = await api_client.get(f"/tasks/{task_id}/analytics")
        analytics = TaskAnalytics(**analytics_res.json())
        
        assert analytics.percentage_complete == 100.0
        assert analytics.time_to_complete_hours is not None