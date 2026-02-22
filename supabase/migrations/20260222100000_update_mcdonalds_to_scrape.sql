-- マクドナルドのRSSフィードが廃止されたため、HTMLスクレイピングに変更
update sources
set
  url  = 'https://www.mcdonalds.co.jp/company/news/',
  type = 'scrape'
where name = 'マクドナルド';
