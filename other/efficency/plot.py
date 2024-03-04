import matplotlib.pyplot as plt
import pandas as pd

plt.rcParams['font.family'] = 'Heiti TC'

# 更新数据
data_updated_chinese = {
    '执行次数': [100, 1000, 10000, 100000],
    '4': [571.711e-3, 6.161681875, 81.2639595, 1097.018320333],
    '8': [595.460542e-3, 6.594456792, 87.299992875, 1142.851574708],
    '12': [654.573833e-3, 7.078005, 95.284515709, 1168.960288542],
    '16': [719.704792e-3, 7.677571191, 102.285363041, 1223.367387959]
}

df_updated_chinese = pd.DataFrame(data_updated_chinese)

# 计算每次请求的平均执行时间（秒）
for node in ['4', '8', '12', '16']:
    df_updated_chinese[f'{node}_平均执行时间'] = df_updated_chinese[node] / df_updated_chinese['执行次数']

# 设置不同的线型和标记样式
line_styles = {
    '4': {'linestyle': '-', 'marker': 'o'},
    '8': {'linestyle': '--', 'marker': 's'},
    '12': {'linestyle': '-.', 'marker': '^'},
    '16': {'linestyle': ':', 'marker': 'd'}
}

plt.figure(figsize=(10, 6))

# 绘制折线图，展示平均执行时间
for node, style in line_styles.items():
    plt.plot(df_updated_chinese['执行次数'], df_updated_chinese[f'{node}_平均执行时间'], label=f'{node} 节点',
             linestyle=style['linestyle'], marker=style['marker'])

plt.title('链下计算存储系统响应性能测试')
plt.xlabel('执行次数')
plt.ylabel('平均执行时间（秒）')
plt.xscale('log')
plt.yscale('log')
plt.legend()
plt.grid(True, which="both", ls="--")

# 展示图表
plt.show()
