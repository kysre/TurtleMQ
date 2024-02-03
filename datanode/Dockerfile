FROM python:3.10-slim

WORKDIR /usr/src/app

COPY /datanode/requirements.txt ./
RUN pip install -r requirements.txt

COPY /datanode/src ./src

EXPOSE 8000

CMD ["python", "./src/datanode_server.py"]
