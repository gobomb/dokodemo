$(document).ready(function () {

    $("span.glyphicon-off").click(function () {
        $("a#out").hide();
    });
});

// type Info struct {
//     IP         string
//     Status     bool
//     Tuns       []string
//     Ctls       []string
//     TunnelAddr string
//     Domain     string
//     Consnum    int
// }

// 开启服务
//ajax1 = $.getJSON("/status-on", function (r) {
//})

getInfo();


setInterval(function() {
   getInfo();
}, 2000);

// 更新服务器状态
function getInfo() {
    $.getJSON("/info", function (r) {

        if (r['Status'] == true) {
            $("span#server-status").attr('class', 'light green');
            $("span#server-status").attr('status',1);
            $("a#open-server").text("关闭服务器");
        } else {
            $("span#server-status").attr('class', 'light red');
            $("span#server-status").attr('status',0);
            $("a#open-server").text("打开服务器");
        }

        // 更新客户端信息
        arrayCtl = r['Ctls'];
        // console.log(arrayCtl);
        if (arrayCtl != null) {
            for (var key in  arrayCtl) {
                console.log(key);
                $("td#client-id").text(key);
            }
        }

    })
}



// 开关服务器
function changeServer() {
    status=$("span#server-status").attr("status");
    console.log(status);
    if(status==0){
        $.getJSON("/status-on", function (r) {
            alert("on")
        })
    }else{
        $.getJSON("/status-off", function (r) {
            alert("off")
        })
    }

}
// 监听开关服务器器
$(document).ready(function () {
        $("a#open-server").click(function () {
            changeServer()
        })
    }
)
// $("td#server-port").text(r['TunnelAddr']);
// if (r['Status'] == true) {
//     $("span#server-status").attr('class', 'light green');
// } else {
//     $("span#server-status").attr('class', 'light red');
// }

// $("span#server-").attr("light red")
