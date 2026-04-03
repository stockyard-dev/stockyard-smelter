package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-smelter/internal/server";"github.com/stockyard-dev/stockyard-smelter/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9700"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./smelter-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("smelter: %v",err)};defer db.Close();srv:=server.New(db,server.DefaultLimits())
fmt.Printf("\n  Smelter — Self-hosted data transformation API\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("smelter: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
