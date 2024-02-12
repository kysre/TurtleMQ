import sys
import yaml


def get_current_state(docker_compose):
    result = 0
    for service in docker_compose['services']:
        if service.startswith('datanode_'):
            datanode_number = int(service[9:])
            result = max(result, datanode_number + 1)
    return result


def scale_up_compose(docker_compose, state, scale_to):
    datanode_template = """
        image: kysre/turtlemq:datanode-${{DATANODE_IMAGE_TAG}}
        restart: unless-stopped
        depends_on:
            - leader_0
            - leader_1
        environment:
            - HOME_PATH=/var/lib/turtlemq/data/
            - DATANODE_NAME=datanode_{x}
            - DATANODE_PORT=8000
            - LEADER_HOST=leader_0
            - LEADER_PORT=8888
            - PULL_TIMEOUT=10
            - PENDING_TIMEOUT=15
            - CLEANER_PERIOD=3
            - PARTITIONS_COUNT=100
        volumes:
            - datanode_{x}_vol:/var/lib/turtlemq/data/
        deploy:
          resources:
            limits:
              cpus: '0.15'
              memory: 300M
    """
    volume_template = '''
        driver: local
    '''

    for x in range(state, scale_to):
        docker_compose['services'][f'datanode_{x}'] = yaml.safe_load(datanode_template.format(x=x))
        docker_compose['volumes'][f'datanode_{x}_vol'] = yaml.safe_load(volume_template.format(x=x))


def add_datanode_to_prometheus(state, scale_to):
    with open('./prometheus/prometheus.yml', 'a') as f:
        for x in range(state, scale_to):
            f.write(f"\n  - targets: ['datanode_{x}:9000']")


def diff_update_compose():
    docker_compose = yaml.safe_load(open('./docker-compose.yaml', 'r').read())
    state = get_current_state(docker_compose)
    scale_to = int(sys.argv[1])
    if scale_to < state:
        print('Cannot scale down.')
        sys.exit(1)
    scale_up_compose(docker_compose, state, scale_to)
    yaml.dump(docker_compose, open('./docker-compose.yaml', 'w'))
    add_datanode_to_prometheus(state, scale_to)


if __name__ == '__main__':
    diff_update_compose()
