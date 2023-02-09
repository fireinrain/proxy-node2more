var ipv4 = {
  random: function(subnet, mask) {
    let randomIp = Math.floor(Math.random() * Math.pow(2, 32 - mask)) + 1;

    return this.lon2ip(this.ip2lon(subnet) | randomIp);
  },
  ip2lon: function(address) {
    let result = 0;

    address.split('.').forEach(function(octet) {
      result <<= 8;
      result += parseInt(octet, 10);
    });

    return result >>> 0;
  },
  lon2ip: function(lon) {
    return [lon >>> 24, lon >> 16 & 255, lon >> 8 & 255, lon & 255].join('.');
  }
};

function exec() {
  var cdntype = $("#cdntype").val();
  var nodes_num = $("#num").val();
  var out = $("#out");
  var base_url = "https://api.buliang0.cf/cdnyx/?cdntype=";
  var ip_result = [];
  var giplist = [];

  var sample_node = ($("#in").val()).replace(/[\r\n\s]/g, "");
  var vmess_pre = "vmess://";
  var vless_pre = "vless://";
  var trojan_pre = "trojan://";

  if (sample_node.indexOf(vmess_pre) != 0 && sample_node.indexOf(vless_pre) != 0 && sample_node.indexOf(trojan_pre) != 0 ){
    $("#err").html("仅支持vmess、vless和trojan的节点分享链接！");
    return;
  }

  if (cdntype !== "custom") {
    $.ajax({
      url: base_url + cdntype,
      type: "GET",
      async: false,
      success: result=> {
        switch(cdntype)
        {
          case "cf":
            giplist = result.split(/[(\r\n)\r\n]+/);
            break;
          case "cft":
            giplist = result.CLOUDFRONT_GLOBAL_IP_LIST;
            break;
          case "gcore":
            giplist = result.addresses;
            break;
        }
      },
      error:result=>{
        $("#err").html("IP获取失败！网页返回状态码："+result.status);
      }
    });
  } else {
    giplist = $("#cdniplist").val().replace(/^[(\r\n)\r\n\s]+|[(\r\n)\r\n\s]+$/g, '').split(/[(\r\n)\r\n]+/);
  }
  if ($("#getmethod").val() === "queue") {
    for (let index = i = 0; i < nodes_num; i++,index++) {
      if (giplist.length==index) {
        index=0;
      }
      let ip_mask_result = giplist[index].split('/');
      ip_result.push(ipv4.random(ip_mask_result[0], ip_mask_result[1]));
    }
  }else{
    for (let i = 0; i < nodes_num; i++) {
      let index = Math.floor((Math.random() * giplist.length));
      let ip_mask_result = giplist[index].split('/');
      ip_result.push(ipv4.random(ip_mask_result[0], ip_mask_result[1]));
    }
  }
  ip_result = Array.from(new Set(ip_result));
  var nodes = new Array();

  if (sample_node.indexOf(vmess_pre) == 0) {
    node_data = atob(sample_node.slice(vmess_pre.length, sample_node.length))
    var re = /\"add\": ?\"(.*?)\"/;
    var node_host = node_data.match(re)[1];
    node_data = node_data.replace(/(\"host\": ?\")(.*?)(\"\,)/, "$1" + node_host + "$3");

    for (var i = 0; i < ip_result.length; i++) {
      node_data = node_data.replace(/(\"add\": ?\")(.*?)(\"\,)/, "$1" + ip_result[i] + "$3");
      nodes.push(vmess_pre + btoa(node_data) + "\n")
    }
    out.html(nodes);
  } else if (sample_node.indexOf(vless_pre) == 0 | sample_node.indexOf(trojan_pre) == 0) {
    var re = /@(.*?):/;
    var node_host = sample_node.match(re)[1];

    if (sample_node.indexOf("host=") != -1) {
      sample_node = sample_node.replace(/(host=)(.*?)(&)/, "$1" + node_host + "$3");
    } else {
      sample_node = sample_node.replace(/(@)(.*?)(:)(.*?)(\?)/, "$1$2$3$4$5host=" + node_host + "&");
    }
    for (var i = 0; i < ip_result.length; i++) {
      nodes.push(sample_node.replace(/(@)(.*?)(:)/, "$1" + ip_result[i] + "$3") + "\n");
    }
    out.html(nodes);
  } else {
    $("#err").html("仅支持vmess、vless和trojan的节点分享链接！");

  }
}

function isCustom(objSelect, objItemValue) {
  if ($("#cdntype").val() === "custom") {
    document.getElementById("customview").style.display = "";
  } else {
    document.getElementById("customview").style.display = "none";
  }

}

function cleantext() {
  $("#err").html("");
  $("#out").html("");
}