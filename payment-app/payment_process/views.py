import json
import random
import requests
from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
from .models import Invoice

@csrf_exempt
def process_payment(request):
    invoice_number = request.POST.get('invoice_number')
    amount = float(request.POST.get('amount'))

    # Simulasi kegagalan 10%
    success = random.choices([True, False], weights=[90, 10])[0]

    invoice = Invoice.objects.create(
        invoice_number=invoice_number,
        amount=amount,
        success=success,
    )

    # Kirim webhook ke Ticket App
    webhook_url = "WEBHOOK_URL"  # Ganti dengan konfigurasi yang sesuai
    payload = {
        'invoice_number': invoice.invoice_number,
        'amount': invoice.amount,
        'success': invoice.success,
    }
    requests.post(webhook_url, data=json.dumps(payload), headers={'Content-Type': 'application/json'})

    return JsonResponse({'message': 'Payment processed successfully'})

