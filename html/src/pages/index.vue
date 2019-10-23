<template>
    <div>
        <el-row>
            <el-col :span="6">
                <el-card class="box-card">
                    <div slot="header" class="clearfix">
                        <span>数据源.在线</span>
                        <el-tag style="float: right;" size="small" type="success">
                            周
                        </el-tag>
                    </div>
                    <div style="text-align: center">
                        <el-progress type="circle" :percentage="100"></el-progress>
                        <el-divider content-position="center">极光云</el-divider>
                        <el-tag size="small">数据源：6</el-tag>
                    </div>
                </el-card>
            </el-col>
            <el-col :span="6">
                <el-card class="box-card">
                    <div slot="header" class="clearfix">
                        <span>数据应用</span>
                        <el-tag style="float: right;" type="danger" size="small">
                            全部
                        </el-tag>
                    </div>
                    <div style="text-align: center">
                        <el-progress type="circle" :percentage="25"></el-progress>
                        <el-divider content-position="center">极光云</el-divider>
                        <el-tag size="small">数据应用：31</el-tag>
                    </div>
                </el-card>
            </el-col>
            <el-col :span="6">
                <el-card class="box-card">
                    <div slot="header" class="clearfix">
                        <span>容器服务</span>
                        <el-tag style="float: right;" type="danger" size="small">
                            全部
                        </el-tag>
                    </div>
                    <div style="text-align: center">
                        <el-progress type="circle" :percentage="86"></el-progress>
                        <el-divider content-position="center">极光云</el-divider>
                        <el-tag size="small">容器服务：11</el-tag>
                    </div>
                </el-card>
            </el-col>
            <el-col :span="6">
                <el-card class="box-card">
                    <div slot="header" class="clearfix">
                        <span>监控报警</span>
                        <el-tag style="float: right;" size="small">
                            月
                        </el-tag>
                    </div>
                    <div style="text-align: center">
                        <el-progress type="circle" :percentage="3" status="warning"></el-progress>
                        <el-divider content-position="center">极光云</el-divider>
                        <el-tag size="small">报警：3</el-tag>
                    </div>
                </el-card>
            </el-col>
        </el-row>
        <el-row style="margin: 0.5rem 0 0 0">
            <el-col :span="24">
                <el-card class="box-card">
                    <div slot="header" class="clearfix">
                        <span>容器服务</span>
                        <el-tag style="float: right;" type="success" size="small">
                            统计
                        </el-tag>
                    </div>
                    <div style="text-align: center">
                        <div id="c1"></div>
                    </div>
                </el-card>
            </el-col>
            <el-col :span="8">
            </el-col>
        </el-row>
        <el-row style="margin: 0.5rem 0 0 0">
            <el-col :span="24">
                <el-card class="box-card">
                    <div slot="header" class="clearfix">
                        <span>服务依赖</span>
                        <el-tag style="float: right;" type="info" size="small">
                            分析
                        </el-tag>
                    </div>
                    <div style="text-align: center">
                        <div id="c2"></div>
                    </div>
                </el-card>
            </el-col>
            <el-col :span="8">
            </el-col>
        </el-row>
    </div>
</template>

<script>
    import G2 from '@antv/g2';
    import G6 from '@antv/g6'
    export default {
        name: "index",
        methods: {},
        mounted() {
            var data = [{
                type: '分类一',
                value: 27
            }, {
                type: '分类二',
                value: 25
            }, {
                type: '分类三',
                value: 18
            }, {
                type: '分类四',
                value: 15
            }, {
                type: '分类五',
                value: 10
            }, {
                type: 'Other',
                value: 5
            }];

            // 根据比例，获取两点之间的点
            function getPoint(p0, p1, ratio) {
                return {
                    x: (1 - ratio) * p0.x + ratio * p1.x,
                    y: (1 - ratio) * p0.y + ratio * p1.y
                };
            }

            var pointRatio = 0.7; // 设置开始变成圆弧的位置 0.7
            // 可以通过调整这个数值控制分割空白处的间距，0-1 之间的数值
            var sliceNumber = 0.005;

            // 自定义 other 的图形，增加两条线
            G2.Shape.registerShape('interval', 'platelet', {
                draw: function draw(cfg, container) {
                    cfg.points[1].y = cfg.points[1].y - sliceNumber;
                    cfg.points[2].y = cfg.points[2].y - sliceNumber;
                    var centerPoint = {
                        x: cfg.points[3].x,
                        y: (cfg.points[2].y + cfg.points[3].y) / 2
                    };
                    centerPoint = this.parsePoint(centerPoint);
                    var points = this.parsePoints(cfg.points);
                    var path = [];
                    var tmpPoint1 = getPoint(points[0], points[3], pointRatio);
                    var tmpPoint2 = getPoint(points[1], points[2], pointRatio);
                    path.push(['M', points[0].x, points[0].y]);
                    path.push(['L', tmpPoint1.x, tmpPoint1.y]);
                    path.push(['Q', points[3].x, points[3].y, centerPoint.x, centerPoint.y]);
                    path.push(['Q', points[2].x, points[2].y, tmpPoint2.x, tmpPoint2.y]);
                    path.push(['L', points[1].x, points[1].y]);
                    path.push(['z']);
                    return container.addShape('path', {
                        attrs: {
                            fill: cfg.color,
                            path: path
                        }
                    });
                }
            });

            var chart = new G2.Chart({
                container: 'c1',
                forceFit: true,
                height: 260,
                padding: [40, 0]
            });

            chart.source(data);
            chart.coord('theta');
            chart.intervalStack().position('value').color('type').shape('platelet').label('type');

            chart.render();

            var _this = this;

            var ERROR_COLOR = "#F5222D";
            var getNodeConfig = function getNodeConfig(node) {
                if (node.nodeError) {
                    return {
                        basicColor: ERROR_COLOR,
                        fontColor: "#FFF",
                        borderColor: ERROR_COLOR,
                        bgColor: "#E66A6C"
                    };
                }
                var config = {
                    basicColor: "#722ED1",
                    fontColor: "#722ED1",
                    borderColor: "#722ED1",
                    bgColor: "#F6EDFC"
                };
                switch (node.type) {
                    case "root": {
                        config = {
                            basicColor: "#E3E6E8",
                            fontColor: "rgba(0,0,0,0.85)",
                            borderColor: "#E3E6E8",
                            bgColor: "#F7F9FA"
                        };
                        break;
                    }
                    default:
                        break;
                }
                return config;
            };

            var COLLAPSE_ICON = function COLLAPSE_ICON(x, y, r) {
                return [["M", x - r, y], ["a", r, r, 0, 1, 0, r * 2, 0], ["a", r, r, 0, 1, 0, -r * 2, 0], ["M", x - r + 4, y], ["L", x - r + 2 * r - 4, y]];
            };
            var EXPAND_ICON = function EXPAND_ICON(x, y, r) {
                return [["M", x - r, y], ["a", r, r, 0, 1, 0, r * 2, 0], ["a", r, r, 0, 1, 0, -r * 2, 0], ["M", x - r + 4, y], ["L", x - r + 2 * r - 4, y], ["M", x - r + r, y - r + 4], ["L", x, y + r - 4]];
            };

            var nodeBasicMethod = {
                createNodeBox: function createNodeBox(group, config, width, height, isRoot) {
                    /* 最外面的大矩形 */
                    var container = group.addShape("rect", {
                        attrs: {
                            x: 0,
                            y: 0,
                            width: width,
                            height: height
                        }
                    });
                    if (!isRoot) {
                        /* 左边的小圆点 */
                        group.addShape("circle", {
                            attrs: {
                                x: 3,
                                y: height / 2,
                                r: 6,
                                fill: config.basicColor
                            }
                        });
                    }
                    /* 矩形 */
                    group.addShape("rect", {
                        attrs: {
                            x: 3,
                            y: 0,
                            width: width - 19,
                            height: height,
                            fill: config.bgColor,
                            stroke: config.borderColor,
                            radius: 2,
                            cursor: "pointer"
                        }
                    });

                    /* 左边的粗线 */
                    group.addShape("rect", {
                        attrs: {
                            x: 3,
                            y: 0,
                            width: 3,
                            height: height,
                            fill: config.basicColor,
                            radius: 1.5
                        }
                    });
                    return container;
                },
                /* 生成树上的 marker */
                createNodeMarker: function createNodeMarker(group, collapsed, x, y) {
                    group.addShape("circle", {
                        attrs: {
                            x: x,
                            y: y,
                            r: 13,
                            fill: "rgba(47, 84, 235, 0.05)",
                            opacity: 0,
                            zIndex: -2
                        },
                        className: "collapse-icon-bg"
                    });
                    group.addShape("marker", {
                        attrs: {
                            x: x,
                            y: y,
                            radius: 7,
                            symbol: collapsed ? EXPAND_ICON : COLLAPSE_ICON,
                            stroke: "rgba(0,0,0,0.25)",
                            fill: "rgba(0,0,0,0)",
                            lineWidth: 1,
                            cursor: "pointer"
                        },
                        className: "collapse-icon"
                    });
                },
                afterDraw: function afterDraw(cfg, group) {
                    /* 操作 marker 的背景色显示隐藏 */
                    var icon = group.findByClassName("collapse-icon");
                    if (icon) {
                        var bg = group.findByClassName("collapse-icon-bg");
                        icon.on("mouseenter", function() {
                            bg.attr("opacity", 1);
                            graph.get("canvas").draw();
                        });
                        icon.on("mouseleave", function() {
                            bg.attr("opacity", 0);
                            graph.get("canvas").draw();
                        });
                    }
                    /* ip 显示 */
                    var ipBox = group.findByClassName("ip-box");
                    if (ipBox) {
                        /* ip 复制的几个元素 */
                        var ipLine = group.findByClassName("ip-cp-line");
                        var ipBG = group.findByClassName("ip-cp-bg");
                        var ipIcon = group.findByClassName("ip-cp-icon");
                        var ipCPBox = group.findByClassName("ip-cp-box");

                        var onMouseEnter = function onMouseEnter() {
                            _this.ipHideTimer && clearTimeout(_this.ipHideTimer);
                            ipLine.attr("opacity", 1);
                            ipBG.attr("opacity", 1);
                            ipIcon.attr("opacity", 1);
                            graph.get("canvas").draw();
                        };
                        var onMouseLeave = function onMouseLeave() {
                            _this.ipHideTimer = setTimeout(function() {
                                ipLine.attr("opacity", 0);
                                ipBG.attr("opacity", 0);
                                ipIcon.attr("opacity", 0);
                                graph.get("canvas").draw();
                            }, 100);
                        };
                        ipBox.on("mouseenter", function() {
                            onMouseEnter();
                        });
                        ipBox.on("mouseleave", function() {
                            onMouseLeave();
                        });
                        ipCPBox.on("mouseenter", function() {
                            onMouseEnter();
                        });
                        ipCPBox.on("mouseleave", function() {
                            onMouseLeave();
                        });
                        ipCPBox.on("click", function() {});
                    }
                },
                setState: function setState(name, value, item) {
                    var hasOpacityClass = ["ip-cp-line", "ip-cp-bg", "ip-cp-icon", "ip-cp-box", "ip-box", "collapse-icon-bg"];
                    var group = item.getContainer();
                    var childrens = group.get("children");
                    graph.setAutoPaint(false);
                    if (name === "emptiness") {
                        if (value) {
                            childrens.forEach(function(shape) {
                                if (hasOpacityClass.indexOf(shape.get("className")) > -1) {
                                    return;
                                }
                                shape.attr("opacity", 0.4);
                            });
                        } else {
                            childrens.forEach(function(shape) {
                                if (hasOpacityClass.indexOf(shape.get("className")) > -1) {
                                    return;
                                }
                                shape.attr("opacity", 1);
                            });
                        }
                    }
                    graph.setAutoPaint(true);
                }
            };

            /**
             * 计算字符串的长度
             * @param {string} str 指定的字符串
             */
            var calcStrLen = function calcStrLen(str) {
                var len = 0;
                for (var i = 0; i < str.length; i++) {
                    if (str.charCodeAt(i) > 0 && str.charCodeAt(i) < 128) {
                        len++;
                    } else {
                        len += 2;
                    }
                }
                return len;
            };

            G6.registerNode('card-node', {
                drawShape: function drawShape(cfg, group) {
                    var config = getNodeConfig(cfg);
                    var isRoot = cfg.type === "root";
                    var nodeError = cfg.nodeError;
                    /* 最外面的大矩形 */
                    var container = nodeBasicMethod.createNodeBox(group, config, 243, 64, isRoot);

                    if (cfg.type !== "root") {
                        /* 上边的 type */
                        group.addShape("text", {
                            attrs: {
                                text: cfg.type,
                                x: 3,
                                y: -10,
                                fontSize: 12,
                                textAlign: "left",
                                textBaseline: "middle",
                                fill: "rgba(0,0,0,0.65)"
                            }
                        });
                    }

                    var ipWidth = 0;
                    if (cfg.ip) {
                        /* ip start */
                        /* ipBox */
                        var ipRect = group.addShape("rect", {
                            attrs: {
                                fill: nodeError ? null : "#FFF",
                                stroke: nodeError ? "rgba(255,255,255,0.65)" : null,
                                radius: 2,
                                cursor: "pointer"
                            }
                        });

                        /* ip */
                        var ipText = group.addShape("text", {
                            attrs: {
                                text: cfg.ip,
                                x: 0,
                                y: 19,
                                fontSize: 12,
                                textAlign: "left",
                                textBaseline: "middle",
                                fill: nodeError ? "rgba(255,255,255,0.85)" : "rgba(0,0,0,0.65)",
                                cursor: "pointer"
                            }
                        });

                        var ipBBox = ipText.getBBox();
                        /* ip 的文字总是距离右边 12px */
                        ipText.attr({
                            x: 224 - 12 - ipBBox.width
                        });
                        /* ipBox */
                        ipRect.attr({
                            x: 224 - 12 - ipBBox.width - 4,
                            y: ipBBox.minY - 5,
                            width: ipBBox.width + 8,
                            height: ipBBox.height + 10
                        });

                        /* 在 IP 元素上面覆盖一层透明层，方便监听 hover 事件 */
                        group.addShape("rect", {
                            attrs: {
                                stroke: "",
                                cursor: "pointer",
                                x: 224 - 12 - ipBBox.width - 4,
                                y: ipBBox.minY - 5,
                                width: ipBBox.width + 8,
                                height: ipBBox.height + 10,
                                fill: "#fff",
                                opacity: 0
                            },
                            className: "ip-box"
                        });

                        /* copyIpLine */
                        group.addShape("rect", {
                            attrs: {
                                x: 194,
                                y: 7,
                                width: 1,
                                height: 24,
                                fill: "#E3E6E8",
                                opacity: 0
                            },
                            className: "ip-cp-line"
                        });
                        /* copyIpBG */
                        group.addShape("rect", {
                            attrs: {
                                x: 195,
                                y: 8,
                                width: 22,
                                height: 22,
                                fill: "#FFF",
                                cursor: "pointer",
                                opacity: 0
                            },
                            className: "ip-cp-bg"
                        });
                        /* copyIpIcon */
                        group.addShape("image", {
                            attrs: {
                                x: 200,
                                y: 13,
                                height: 12,
                                width: 10,
                                img: "https://os.alipayobjects.com/rmsportal/DFhnQEhHyPjSGYW.png",
                                cursor: "pointer",
                                opacity: 0
                            },
                            className: "ip-cp-icon"
                        });
                        /* 放一个透明的矩形在 icon 区域上，方便监听点击 */
                        group.addShape("rect", {
                            attrs: {
                                x: 195,
                                y: 8,
                                width: 22,
                                height: 22,
                                fill: "#FFF",
                                cursor: "pointer",
                                opacity: 0
                            },
                            className: "ip-cp-box",
                            tooltip: "复制IP"
                        });

                        var ipRectBBox = ipRect.getBBox();
                        ipWidth = ipRectBBox.width;
                        /* ip end */
                    }

                    /* name */
                    var nameText = group.addShape("text", {
                        attrs: {
                            text: cfg.name,
                            x: 19,
                            y: 19,
                            fontSize: 14,
                            fontWeight: 700,
                            textAlign: "left",
                            textBaseline: "middle",
                            fill: config.fontColor,
                            cursor: "pointer"
                            // tooltip: cfg.name,
                        }
                    });

                    /* 下面的文字 */
                    var remarkText = group.addShape("text", {
                        attrs: {
                            text: cfg.keyInfo,
                            x: 19,
                            y: 45,
                            fontSize: 14,
                            textAlign: "left",
                            textBaseline: "middle",
                            fill: config.fontColor,
                            cursor: "pointer"
                        }
                    });

                    if (nodeError) {
                        group.addShape("text", {
                            attrs: {
                                x: 191,
                                y: 62,
                                text: '⚠️',
                                fill: '#000',
                                fontSize: 18
                            }
                        });
                    }

                    var hasChildren = cfg.children && cfg.children.length > 0;
                    if (hasChildren) {
                        nodeBasicMethod.createNodeMarker(group, cfg.collapsed, 236, 32);
                    }
                    return container;
                },
                afterDraw: nodeBasicMethod.afterDraw,
                setState: nodeBasicMethod.setState
            }, "single-shape");

            var data = {
                nodes: [{
                    name: 'cardNodeApp',
                    ip: '127.0.0.1',
                    nodeError: true,
                    type: 'root',
                    keyInfo: 'this is a card node info',
                    x: 100,
                    y: 50
                }, {
                    name: 'cardNodeApp',
                    ip: '127.0.0.1',
                    nodeError: false,
                    type: 'subRoot',
                    keyInfo: 'this is sub root',
                    x: 100,
                    y: 150
                }, {
                    name: 'cardNodeApp',
                    ip: '127.0.0.1',
                    nodeError: false,
                    type: 'subRoot',
                    keyInfo: 'this is sub root',
                    x: 100,
                    y: 250,
                    children: [{
                        name: 'sub'
                    }]
                }],
                edges: []
            };

            var graph = new G6.Graph({
                container: 'c2',
                width: 400,
                height: 400,
                defaultNode: {
                    shape: 'card-node'
                }
            });

            graph.data(data);
            graph.render();
        }
    }
</script>

<style scoped>
    .box-card {
        margin: 0 0.5rem;
    }
</style>