table = [{"id": 1, "text": 'ab'},
         {"id": 2, "text": 'a'}]
    
cur_str = 'ab'
        
def get_level(str):
    level = 0
    
    
    for elem in table:
        if str.find(elem["text"]) != -1 and str != elem["text"]:
            level = max(level, get_level(elem["text"]) + 1)
    
    return level

print(get_level(cur_str))
print(get_level('a'))
