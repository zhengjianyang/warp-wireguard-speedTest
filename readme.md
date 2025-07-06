# WireGuard-Ping

这是一个使用Go语言编写的高性能工具，用于并发扫描指定的IP地址段和端口，测试其到WireGuard节点的延迟，并找出延迟最低的节点。

## 功能特性

- **高性能并发扫描**: 基于 [ants](https://github.com/panjf2000/ants) 协程池，能够以极高的效率并发扫描大量IP地址和端口，充分利用系统资源。
- **灵活配置**: 通过 `config.yaml` 文件可以轻松配置扫描的IP段（CIDR格式）、端口、并发数等。
- **延迟测试**: 对每个 `IP:端口` 组合进行多次ping测试，确保结果的准确性。
- **结果排序**: 对有效的节点按延迟从小到大进行排序。
- **结果保存**: 将所有成功的测试结果保存到CSV文件中，便于后续分析。
- **Top 10展示**: 在控制台直接输出延迟最低的前10个节点。

## 环境要求

- [Go](https://golang.org/) (1.22 或更高版本)

**重要提示**: 必须使用 Go 1.22 或更高版本进行编译。本项目利用了 Go 1.22 中修复的循环变量捕获问题，以确保并发任务的正确性。

## 安装与使用

1.  **克隆项目**
    ```bash
    git clone <your-repo-url>
    cd wireguardPing
    ```

2.  **安装依赖**
    ```bash
    go mod tidy
    ```

3.  **编译项目**
    ```bash
    go build
    ```

4.  **修改配置**
    根据您的需求修改 `config.yaml` 文件。

5.  **运行程序**
    ```bash
    ./wireguardPing
    ```

## 配置文件

程序的所有配置都在 `config.yaml` 文件中。

```yaml
# config.yaml

# 保存结果文件名
saveFileName: "results.csv"

# 最大并发协程数
threads: 200

# 对每个IP:Port组合的成功测试次数要求
pingCount: 3

# 要测试的端口列表
ports:
  - 500
  - 1701
  - 2408
  - 4500

# 要扫描的IP地址段 (CIDR格式)
ipv4CIDRs:
  - "162.159.192.0/24"
  - "162.159.193.0/24"
  - "188.114.96.0/24"
  - "188.114.97.0/24"
```

- `saveFileName`: 保存扫描结果的CSV文件名。
- `threads`: 并发执行扫描任务的协程（goroutine）数量。
- `pingCount`: 每个节点需要成功连接的次数才被认为是有效节点。
- `ports`: 需要扫描的端口列表。
- `ipv4CIDRs`: 需要扫描的IP地址范围，使用CIDR格式。

## 输出说明

程序运行结束后，会执行以下操作：

1.  **控制台输出**: 打印出延迟最低的前10个 `IP:端口` 和对应的延迟时间。
    ```
    前10条延迟最低的结果>>>
    [1] IP: 188.114.99.238, Port: 1701, Latency: 171ms
    [2] IP: 188.114.99.101, Port: 500, Latency: 172ms
    ...
    ```

2.  **CSV文件**: 将所有扫描到的有效节点及其延迟信息保存到 `config.yaml` 中 `saveFileName` 指定的文件中（默认为 `results.csv`）。

    **`results.csv` 文件示例:**
    ```csv
    IP,Port,Latency(ms)
    188.114.99.238,1701,171.00
    188.114.99.101,500,172.00
    188.114.99.210,500,172.00
    ...
    ```

## 致谢

本项目的灵感和部分实现参考了 [CloudflareWarpSpeedTest](https://github.com/peanut996/CloudflareWarpSpeedTest) 项目，特此感谢。

## 许可证

本项目采用 [MIT](LICENSE) 许可证。