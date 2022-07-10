create keyspace if not exists ctrlplane 
  with replication = {'class': 'SimpleStrategy', 'replication_factor': 1};