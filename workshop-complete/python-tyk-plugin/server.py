from concurrent import futures
import grpc, time, json

import pika

import coprocess_object_pb2_grpc as coprocess
from coprocess_object_pb2 import Object

AMQP_URL = 'amqp://server'

class Todo():
  def __init__(self):
    pass

class Dispatcher(coprocess.DispatcherServicer):
  def Dispatch(self, obj, context):
    if obj.hook_name != 'TodoRabbitHook':
      print("Unknown hook")
      return None
    req = obj.request
    path = obj.request.url.replace('/todos', '', -1)
    id_str = path.replace('/', '', 1)

    routing_key = ''
    todo = {}

    if req.method == 'GET':
      if path == '/':
        routing_key = 'index'
        todo['user'] = obj.session.alias
      else:
        routing_key = 'show'
        todo['id'] = id_str
    elif req.method == 'POST':
      routing_key = 'store'
      todo = json.loads(req.raw_body)
    elif req.method == 'DELETE':
      routing_key = 'delete'
      todo['id'] = id_str
    elif req.method == 'PATCH':
      routing_key = 'update'
      todo['id'] = id_str
      todo = json.loads(req.raw_body)
    else:
      return None
    
    todo['user'] = obj.session.alias

    # Handle AMQP interactions:
    messages = self.handle_rpc(routing_key, todo)

    # Override the response:
    obj.request.return_overrides.response_code = 200
    obj.request.return_overrides.response_error = json.dumps(messages)
    obj.request.return_overrides.headers["Content-Type"] = "application/json"
    return obj
    
  def handle_rpc(self, routing_key, todo):
    todo_json = json.dumps(todo)
    params = pika.URLParameters(AMQP_URL)
    connection = pika.BlockingConnection(params)
    channel = connection.channel()
    channel.queue_declare(queue="reply-to",
                              exclusive=True,
                              auto_delete=True)
    channel.basic_publish(exchange='todos',
                      routing_key=routing_key,
                      body=todo_json,
                      properties=pika.BasicProperties(reply_to='reply-to'))
    channel = connection.channel()
    messages = []

    # Loop until we get a message:
    while True:
      method_frame, header_frame, body = channel.basic_get(queue='reply-to')
      if len(messages) > 0:
        break
      if body is not None:
        messages = json.loads(body)
    connection.close()
    return messages




  def DispatchEvent(self, request, context):
    print('Dispatch is called')
    object = Object()
    return object

server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
coprocess.add_DispatcherServicer_to_server(Dispatcher(), server)
print('Starting server')
server.add_insecure_port('127.0.0.1:9000')
server.start()

try:
  while True:
    time.sleep(86400)
except KeyboardInterrupt:
  print('Stopping server')
