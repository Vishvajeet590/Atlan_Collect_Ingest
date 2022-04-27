# Atlan_Collect_Ingest

Steps to run Ingest Server

1. Clone this repo

        git clone https://github.com/Vishvajeet590/Atlan_Collect_Ingest.git

2. Build the image

        cd ..
         docker build -t ingest -f Dockerfile.ingest ./

    This will take a few minutes.
 
3. Before runnig this image be sure the you have started the Collect Server Image from [the Atlan_Collect_Challenge repository](https://github.com/Vishvajeet590/Atlan_Collect_Challenge) and  and ensure Postgres DBs and RabbitMq is running.
4. Run the image's default command, which should start everything up. The `-p` option forwards the container's port 8080 to port 8000 on the host aong with the postgres DB url as env variable. (Note that the host will actually be a guest if you are using boot2docker, so you may need to re-forward the port in VirtualBox.)

        docker run --rm -p 8081:8080 -e DATABASE_URL='postgres://vishwajeet:docker@localhost:5432/Hermes?&pool_max_conns=10'  --network="host" ingest

    
    
