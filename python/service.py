import cherrypy
import json


class formater(object):
    @cherrypy.expose
    def index(self):
        cl = cherrypy.request.headers['Content-Length']
        raw_message = False
        rawbody = cherrypy.request.body.read(int(cl))
        body = json.loads(rawbody)
        if ('message' in body):
            try:
                raw_message = json.loads(body['message'])
            except ValueError:
                raw_message = False     
        # if request is json write to json text file
        if (raw_message):
            f = open("json_raw.txt", "a")

            f.write(json.dumps(raw_message)+ "\n")

            return json.dumps({"log_type": "json", "log_file": "json_raw.txt"})
        # else request write to text file    
        else :         
            f = open("string_raw.txt", "a")
            f.write(body['message']+ "\n")
            return json.dumps({"log_type": "string", "log_file": "string_raw.txt"})    

cherrypy.server.socket_host = '0.0.0.0'
cherrypy.quickstart(formater())
