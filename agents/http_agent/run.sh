# docker run --rm --name agent http_agent \
# -e DASHBOARD_URL="http://localhost:8080/dashboard/v1/register" \
# -e TARGET_HOST="google.com:80"

./http_agent -httpAddr=7066 -dashboardURL=http://localhost:8080/dashboard/v1/register -targetHost=github.com:443
