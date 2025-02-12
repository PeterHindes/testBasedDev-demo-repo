from datetime import datetime
from pydantic import BaseModel, conint, constr
from typing import Optional

class TaskCreate(BaseModel):
    name: constr(min_length=3, max_length=50)
    start_date: datetime
    due_date: datetime
    priority: conint(ge=1, le=5)

class TaskResponse(BaseModel):
    id: str
    name: str
    start_date: datetime
    due_date: datetime
    priority: int
    status: str
    created_at: datetime
    completed_at: Optional[datetime] = None

class TaskAnalytics(BaseModel):
    days_until_due: int
    completion_health: str
    percentage_complete: float
    time_to_complete_hours: Optional[float] = None

class SystemStats(BaseModel):
    avg_completion_hours: Optional[float] = None
    completion_rate: float
    tasks_by_priority: dict[int, int]