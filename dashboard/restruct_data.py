import pandas as pd

df = pd.read_excel("22fw_주문내역.xlsx")

res = pd.DataFrame({
    "구분": [],
    "시즌": [],
    "상품코드": [],
    "SKU": [],
    "상품명": [],
    "칼라": [],
    "사이즈": [],
    "원가": [],
    "TAG가": [],
    "지점": [],
    "날짜-주차": [],
    "날짜": [],
    "수량": [],
})
    
for i in range(1, len(df)):
     for j in range(7, len(df.columns)):
        if df.iloc[i][j] > 0:
            data = {
                "구분": "판매",
                "시즌": "22FW",
                "상품코드": df.iloc[i][1],
                "SKU": df.iloc[i][1] + "-" + df.iloc[i][3] + "-" + df.iloc[i][4],
                "상품명": df.iloc[i][2],
                "칼라": df.iloc[i][3],
                "사이즈": df.iloc[i][4],
                "원가": df.iloc[i][5],
                "TAG가": df.iloc[i][6],
                "지점": df.columns[j].split('.')[0],
                "날짜-주차": df.iloc[0][j],
                "날짜": "2022/" + df.iloc[0][j].split("~")[0],
                "수량": df.iloc[i][j]
            }
            print(data)
            res = pd.concat([res, pd.DataFrame([data])])
        if df.iloc[i][j] < 0:
            data = {
                "구분": "환불",
                "시즌": "22FW",
                "상품코드": df.iloc[i][1],
                "SKU": df.iloc[i][1] + "-" + df.iloc[i][3] + "-" + df.iloc[i][4],
                "상품명": df.iloc[i][2],
                "칼라": df.iloc[i][3],
                "사이즈": df.iloc[i][4],
                "원가": df.iloc[i][5],
                "TAG가": df.iloc[i][6],
                "지점": df.columns[j].split('.')[0],
                "날짜-주차": df.iloc[0][j],
                "날짜": "2022/" + df.iloc[0][j].split("~")[0],
                "수량": df.iloc[i][j]
            }
            res = pd.concat([res, pd.DataFrame([data])])
            
# 엑셀 파일 출력
output_excel_path = "22fw_출력결과.xlsx"
res.to_excel(output_excel_path, index=False)

print("엑셀 파일로 출력이 완료되었습니다.")
        
