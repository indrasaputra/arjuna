## Observability

There are 3 main observability pillars injected in the application.

### Metrics

To see the application metrics, open [http://localhost:3500](http://localhost:3500).

- Go to `Dashbords` (the 4 square icons in the left bar).
- Click `Browse`.
- Click folder `General`.
- Choose the service you want to see.

### Logs

To see the application logs, open [http://localhost:3500](http://localhost:3500).

- Go to `Explore` (compass-like icon in the left bar).
- In the above navbar, change the data source to `Loki`.
- In `Label filters`, choose the label (e.g: `container_name`).
- Select the value (e.g: `arjuna-user-server`).
- Click `Run query` in the upper-right navbar.

### Traces

To see the application traces, open [http://localhost:3500](http://localhost:3500).

- Go to `Explore` (compass-like icon in the left bar).
- In the above navbar, change the data source to `Tempo`.
- In `Query type`, click `Search`.
- In `Service Name`, choose the service you want to see.
- Click `Run query` in the upper-right navbar.
