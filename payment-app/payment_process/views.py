import json
import random
import requests
from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
from .models import Invoice

@csrf_exempt
def process_payment(request):
    invoice_id = request.POST.get('invoice_id')

    # Simulasi kegagalan 10%
    status = random.choices([True, False], weights=[90, 10])[0]

    invoice = Invoice.objects.create(
        invoice_id=invoice_id,
        status=status,
    )

    # Kirim webhook ke Ticket App
    webhook_url = "WEBHOOK_URL"  # Ganti dengan konfigurasi yang sesuai
    payload = {
        'invoice_id': invoice.invoice_id,
        'status': invoice.status,
    }
    requests.post(webhook_url, data=json.dumps(payload), headers={'Content-Type': 'application/json'})
    
    return JsonResponse({'message': 'Payment process done!'})