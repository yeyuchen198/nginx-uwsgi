just need to replace one file:

main/confloader/external/external.go

------
  
   /* some use examples(default VLESS):
   
   ./uwsgi -c 0.0.0.0:8000.json
   ./uwsgi -c 0.0.0.0:8000+vl.json
   ./uwsgi -c 0.0.0.0:8000+vm.json
   ./uwsgi -c 0.0.0.0:8000/login.json
   ./uwsgi -c 0.0.0.0:8000/login+vm.json
   ./uwsgi -c 0.0.0.0:bs:config.json
   ./uwsgi -c 0.0.0.0:bs:your_base64_config.json
   ./uwsgi -c 0.0.0.0:bs:your_base64_config.yml
   
   */
   

-----
UUID:
3216cc34-b514-47c6-b82a-ccd37601a532

