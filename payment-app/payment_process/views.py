import json
import random
import requests
import asyncio
from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
from django.urls import reverse
from .models import Invoice
from asgiref.sync import async_to_sync, sync_to_async
import aiohttp
import threading
import time

@csrf_exempt
def hello_world(request):
    url = "http://app:3000/" 

    try:
        response = requests.get(url)
        response.raise_for_status()

        try:
            data = response.json() 
            return JsonResponse(data)
        except ValueError:
            print("Raw Response Content:", response.content.decode())
            return JsonResponse({'error': 'Invalid JSON data in response'}, status=500)
            
    except requests.RequestException as e:
        return JsonResponse({'error': str(e)}, status=500)

def generate_unique_invoice_id():
    # Fungsi ini harus memastikan invoice_id adalah unik dan belum ada di database
    while True:
        invoice_id = generate_random_invoice_id()
        if not Invoice.objects.filter(invoice_id=invoice_id).exists():
            return invoice_id

def generate_random_invoice_id():
    # Implementasi sederhana untuk menghasilkan invoice_id yang unik
    return ''.join(random.choices('ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789', k=8))

def process_payment(request):
    invoice_id = generate_random_invoice_id()
    
    # Simulate 10% failure
    status = random.choices([True, False], weights=[90, 10])[0]

    invoice = Invoice.objects.create(
        invoice_id=invoice_id,
        status=status,
    )

    payment_url = request.build_absolute_uri(reverse('process_payment'))

    response_data = {
        'invoice_id': invoice_id,
        'payment_url': payment_url
    }

    webhook_thread = threading.Thread(target=send_webhook, args=(invoice,))
    webhook_thread.start()

    return JsonResponse(response_data)

def send_webhook(invoice):
    time.sleep(10) 
    webhook_url = "http://app:3000/webhook/payment"
    payload = {
        'invoice_id': invoice.invoice_id,
        'status': str(invoice.status),
    }

    try:
        response = requests.post(webhook_url, data=json.dumps(payload), headers={'Content-Type': 'application/json'})
        if response.status_code == 200:
            print("Webhook sent successfully")
        else:
            error_message = response.json() if response.headers.get('content-type') == 'application/json' else response.text
            print(f"Failed to send webhook. Status: {response.status_code}. Error: {error_message}")
    except requests.RequestException as e:
        print(f"Failed to send webhook. Error: {str(e)}")