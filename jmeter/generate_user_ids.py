import csv

# 设置需要生成的用户ID数量
num_user_ids = 200

# 打开文件 user_ids.csv 并写入数据
with open('user_ids.csv', mode='w', newline='') as file:
    writer = csv.writer(file)
    
    # 写入每个userId
    for user_id in range(1, num_user_ids + 1):
        writer.writerow([user_id])

print(f'{num_user_ids} userIds have been written to user_ids.csv')
