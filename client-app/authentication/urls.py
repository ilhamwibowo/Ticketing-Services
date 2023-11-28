from django.urls import path, include
from .views import register_user

urlpatterns = [
    path('register/', register_user, name='register'),
    path("", include("django.contrib.auth.urls")),
]
