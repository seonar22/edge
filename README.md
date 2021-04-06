# edge

Packetframe anycast container orchestration service

### API Reference:

#### Create a container

`POST /containers`

Required parameters:

| Name     | Description                          | Validation          |
|----------|--------------------------------------|---------------------|
| image    | Docker image path                    |                     |
| ports    | List of port map entries ("80:8080") | "uint16:uint16"     |
| env      | Environment map (string to string)   |                     |

Example response (HTTP 200):

```json
{
  "message": "Container created"
}
```

### Delete a container

`DELETE /containers/:container`

Required parameters: None

Example response (HTTP 200):

```json
{
  "message": "Container deleted"
}
```
