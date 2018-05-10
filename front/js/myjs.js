console.log($("#clients>tr>td#client-id").text())



function appendClient(key) {
    var clientId;
    clientId=arrayCtl[key];
    str = "  <tr id='client' new=true>\n" +
        "                <td id='client-id'>"+clientId+
        "</td>\n" +
        "                <td id='tunnel-number'>2</td>\n" +
        "                <td>\n" +
        "                    <div class='dropdown line-div'>\n" +
        "                        <button class='btn btn-default dropdown-toggle' type='button' id='dropdownMenu1'\n" +
        "                                data-toggle='dropdown'>\n" +
        "                            操作\n" +
        "                            <span class='caret'></span>\n" +
        "                        </button>\n" +
        "                        <ul class='dropdown-menu' aria-labelledby='dropdownMenu1'>\n" +
        "                            <li><a href='#'>关闭控制连接</a></li>\n" +
        "                            <li role='presentation'><a href='' role='button' data-toggle='modal' data-target='#myModal'>查看隧道信息</a>\n" +
        "                            </li>\n" +
        "                        </ul>\n" +
        "                    </div>\n" +
        "                </td>\n" +
        "                <td>\n" +
        "                    <span class='light red' id='client-status'></span>\n" +
        "                </td>\n" +
        "            </tr>"

    $("#clients").append(str)
}


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


getInfo();


setInterval(function () {
    getInfo();
}, 2000);

Ctls=new Array();

// 更新服务器状态
function getInfo() {
    $.getJSON("/info", function (r) {

        if (r['Status'] == true) {
            $("span#server-status").attr('class', 'light green');
            $("span#server-status").attr('status', 1);
            $("a#open-server").text("关闭服务器");
        } else {
            $("span#server-status").attr('class', 'light red');
            $("span#server-status").attr('status', 0);
            $("a#open-server").text("打开服务器");
        }

        // 更新客户端信息
        arrayCtl = r['Ctls'];

        if (arrayCtl != null) {
            for (var key in  arrayCtl) {
                i=Ctls.indexOf(arrayCtl[key])
                // console.log(i);
                if (i==-1){
                    Ctls.push(arrayCtl[key])
                    appendClient(key)
                }else{

                }
                ;
            }
        }

    })
}


// 开关服务器
function changeServer() {
    status = $("span#server-status").attr("status");
    console.log(status);
    if (status == 0) {
        $.getJSON("/status-on", function (r) {
            alert("on")
        })
    } else {
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
