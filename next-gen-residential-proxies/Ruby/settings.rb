def get_setting(name)
  get_settings[name]
end

def get_settings
  {
    Username: "",
    Password: "",
    Timeout: 5,
    RequestsRate: 10,
    RetriesNum: 3,
    UrlListName: "url_list_ngrp.txt",
    ProxyAddress: "unblock.oxylabs.io:60000",
  }
end
