from flask import Flask
app = Flask(__name__)


@app.route('/')
def hello():
    return "<h1 style='color:green'>Greetings!</h1>"


@app.route('/<name>')
def hello_name(name):
    return "<h1 style='color:green'>Greetings {}!</h1>".format(name)

if __name__ == '__main__':
    app.run(host='0.0.0.0')
