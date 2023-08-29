require("lru")

print("oop lru")
local cache = lru.LRUCache(5)

cache.add("a", 1)
cache.add("b", 2)
print("len:", cache.size())
cache.add("c", "3a")
cache.add("d", 4)
cache.add("e", 5)
cache.add("f", 6)

print(cache.find("a"))
print(cache.find("d"))
print(cache.find("b"))

print("len:", cache.size())

print("\n遍历cache")
for k in cache.iterator() do
  print(k)
end

print("\n\nmetatable lru")

local a = lru.lrucache(5)

a.a = 1
a.b = 2
print("len:", #a)
a.c = "3a"
a.d = 4
a.e = 5
a.f = 6

print(a.a)
print(a.d)
print(a.b)

print("len:", #a)

print("\n遍历cache")
for k in pairs(a) do
  print(k)
end

local b = lru.lrucache()
print("\n遍历cache", #b)
for k in pairs(b) do
  print(k)
end

print("len:", #b)
