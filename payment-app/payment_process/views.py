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

async def send_webhook(invoice):
    await asyncio.sleep(10)
    webhook_url = "http://app:3000/webhook/payment"  # Replace with the appropriate URL
    payload = {
        'invoice_id': invoice.invoice_id,
        'status': str(invoice.status),
    }

    async with aiohttp.ClientSession() as session:
        async with session.post(webhook_url, data=json.dumps(payload), headers={'Content-Type': 'application/json'}) as response:
            if response.status == 200:
                print("Webhook sent successfully")
            else:
                try:
                    error_message = await response.json()
                    print(f"Failed to send webhook. Status: {response.status}. Error: {error_message}")
                except aiohttp.ContentTypeError:
                    print(f"Failed to send webhook. Status: {response.status}. Unable to read error message.")

@csrf_exempt
def process_payment(request):
    invoice_id = "INV132"
    print("anjrit")
    
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

    asyncio.run(send_webhook(invoice))

    return JsonResponse(response_data)