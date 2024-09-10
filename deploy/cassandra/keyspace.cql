-- NOTE: For production, the default values for class will be 'NetworkTopology' and replicationFactor will be 3.
create keyspace if not exists ctrlplane
  with replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
