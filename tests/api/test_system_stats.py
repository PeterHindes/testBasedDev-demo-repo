import pytest
from httpx import AsyncClient
from tests.models.task_models import SystemStats

@pytest.mark.asyncio
class TestSystemAnalytics:
    async def test_system_stats(self, api_client: AsyncClient, frozen_time):
        # Create mix of tasks
        await api_client.post("/tasks", json={
            "name": "Task 1",
            "start_date": frozen_time.isoformat(),
            "due_date": (frozen_time + timedelta(days=1)).isoformat(),
            "priority": 1
        })
        
        completed_task = await api_client.post("/tasks", json={
            "name": "Task 2",
            "start_date": frozen_time.isoformat(),
            "due_date": (frozen_time + timedelta(hours=2)).isoformat(),
            "priority": 2
        })
        await api_client.put(f"/tasks/{completed_task.json()['id']}/complete")
        
        # Get system stats
        stats_res = await api_client.get("/analytics")
        assert stats_res.status_code == 200
        stats = SystemStats(**stats_res.json()["system_stats"])
        
        assert stats.completion_rate == 50.0
        assert stats.avg_completion_hours is not None
        assert stats.tasks_by_priority == {1: 1, 2: 1}