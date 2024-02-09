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
        environment:
            - HOME_PATH=/var/lib/turtlemq/data/
            - DATANODE_PORT=8000
        deploy:
            mode: replicated
            replicas: ${{DATANODE_REPLICA_COUNT}}
        volumes:
            - datanode_vol_{x}:/var/lib/turtlemq/data/
    """
    volume_template = """
        driver: local
    """

    for x in range(state, scale_to):
        docker_compose['services'][f'datanode_{x}'] = yaml.safe_load(datanode_template.format(x=x))
        docker_compose['volumes'][f'datanode_vol_{x}'] = yaml.safe_load(volume_template.format(x=x))

def diff_update_compose():
    docker_compose = yaml.safe_load(open('./generated.docker-compose.yaml', 'r').read())
    state = get_current_state(docker_compose)
    scale_to = int(sys.argv[1])
    if scale_to < state:
        print('Cannot scale down.')
        sys.exit(1)
    scale_up_compose(docker_compose, state, scale_to)
    yaml.dump(docker_compose, open('./generated.docker-compose.yaml', 'w'))

if __name__ == '__main__':
    diff_update_compose()