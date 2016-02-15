local intsize=2^52 -- dont change
-- technically this could be 2^53 which is the max integer you can have without losing precision
-- but 53 is a prime so its impossible to split up in a reasonable way for multipication
-- plus, 52/2 = 26 bits which is the maximum bits you can multiply together to not overflow 2^53
-- plus plus, 52 is divisable by 4 so conversion to hex is much faster

local bn = {}
local copy

local mt={
	__add=function(a,b)
		a,b=bn(a,true),bn(b,true)
		local o=bn(0)
		if a.neg then
			if b.neg then
				a=-a
				b=-b
				o.neg=true
			else
				return -((-a)-b)
			end
		elseif b.neg then
			return -((-b)-a)
		end
		local o=bn(0)
		local c=0
		for l1=1,math.max(#a,#b) do
			local an,bn=a[l1] or 0,b[l1] or 0
			if an==0 and bn==0 and c==0 and l1==math.max(#a,#b) then
				break
			end
			if c+an+bn>=intsize then
				o[l1]=(an-intsize)+bn+c
				-- subtract intsize before adding second number so it doesnt lose precision
				-- unlike a certain bignum that sacrifices chunk size
				c=1
			else
				o[l1]=an+bn+c
				c=0
			end
		end
		if c>0 then
			o[#o+1]=c
		end
		return o
	end,
	__sub=function(a,b)
		a,b=bn(a,true),bn(b,true)
		local o=bn(0)
		if a==b then
			return o
		end
		if a.neg then
			if b.neg then
				return (-b)-(-a)
			else
				return -((-a)+b)
			end
		elseif b.neg then
			return a+(-b)
		elseif b>a then
			return -(b-a)
		end
		local carry=0
		for l1=1,math.max(#a,#b) do
			local ca,cb=a[l1],b[l1]
			if not cb then
				o[l1]=cb-carry
			else
				o[l1]=ca-cb-carry
			end
			if o[l1]<0 then
				o[l1]=o[l1]+intsize
				carry=1
			else
				carry=0
			end
		end
		assert(carry==0)
		for l1=#o,1,-1 do
			if o[l1]==0 then
				o[l1]=nil
			else
				break
			end
		end
		return o
	end,
	__mul=function(a,b)
		a,b=bn(a,true),bn(b,true)
		if (a.neg and true or false)~=(b.neg and true or false) then
			local o=a*(-b)
			o.neg=true
			return o
		end
		local ans={}
		for l1=1,#a do
			local c=a[l1]
			table.insert(ans,c%(2^26))
			table.insert(ans,math.floor(c/(2^26)))
		end
		if ans[#ans]==0 then
			table.remove(ans)
		end
		local bns={}
		for l1=1,#b do
			local c=b[l1]
			table.insert(bns,c%(2^26))
			table.insert(bns,math.floor(c/(2^26)))
		end
		if bns[#bns]==0 then
			table.remove(bns)
		end
		local o={}
		local carry=0
		for l1=1,#ans do
			for l2=1,#bns do
				carry=ans[l1]*bns[l2]+carry
				carry=carry+(o[((l1-1)+(l2-1))+1] or 0)
				o[((l1-1)+(l2-1))+1]=carry%(2^26)
				carry=math.floor(carry/(2^26))
			end
			if carry>0 then
				o[l1+#bns]=carry
			end
			carry=0
		end
		local t=bn(0)
		t.neg=a.neg
		for l1=1,math.ceil(#o/2) do
			t[l1]=o[((l1-1)*2)+1]+((o[((l1-1)*2)+2] or 0)*(2^26))
		end
		return t
	end,
	__pow=function(a,b)
		a,b=bn(a),bn(b)
		if b.neg then
			return bn(0)
		end
		local result=bn(1)
		if b==0 then
			return result
		end
		while true do
			if b[1]%2==0 then
				b=bn.brshift(b,1)
			else
				b=bn.brshift(b,1)
				result=result*a
				if b==bn(0) then
					return result
				end
			end
			a=a*a
		end
		return result
	end,
	__mod=function(a,b)
		a,b=bn(a,true),bn(b,true)
	end,
	__unm=function(a)
		local o=bn(a)
		o.neg=(not o.neg) and o~=0
		return o
	end,
	__gt=function(a,b)
		a,b=bn(a,true),bn(b,true)
		if a.neg then
			if b.neg then
				return -a<-b
			else
				return false
			end
		elseif b.neg then
			return true
		end
		for l1=math.max(#a,#b),1,-1 do
			local av,bv=a[l1],b[l1]
			if not av then
				return false
			elseif not bv then
				return true
			end
			if av>bv then
				return true
			elseif av<bv then
				return false
			end
		end
		return false
	end,
	__lt=function(a,b)
		a,b=bn(a,true),bn(b,true)
		if a.neg then
			if b.neg then
				return -a>-b
			else
				return true
			end
		elseif b.neg then
			return false
		end
		for l1=math.max(#a,#b),1,-1 do
			local av,bv=a[l1],b[l1]
			if not av then
				return true
			elseif not bv then
				return false
			end
			if av>bv then
				return false
			elseif av<bv then
				return true
			end
		end
		return false
	end,
	__eq=function(a,b)
		a,b=bn(a,true),bn(b,true)
		if (a.neg and true or false)~=(b.neg and true or false) then
			return false
		end
		for l1=math.max(#a,#b),1,-1 do
			if (a[l1] or 0)~=(b[l1] or 0) then
				return false
			end
		end
		return true
	end,
	__tostring=function(a)
		return bn.tostring(a)
	end
}

function copy(n)
	local o=setmetatable({neg=n.neg},mt)
	for l1=1,#n do
		o[l1]=n[l1]
	end
	return o
end

local bnmt={}
setmetatable(bn,bnmt)

function bnmt.__call(s,num,pass)
	if type(num)=="number" then
		num=math.floor(num)
		local o={}
		if num<0 then
			o.neg=true
			num=-num
		end
		if num>=intsize then
			while num>0 do
				table.insert(o,num%intsize)
				num=math.floor(num/intsize)
			end
		else
			o[1]=num
		end
		setmetatable(o,mt)
		return o
	elseif type(num)=="table" then
		return pass and num or copy(num)
	elseif type(num)=="string" then
		return bn.tonumber(num)
	elseif num==nil then
		return setmetatable({},mt)
	else
		error("Unexpected type "..type(num))
	end
end

function bn.brshift(a,b)
	a,b=bn(a),bn(b,true)
	if #b>1 then
		return bn(0)
		-- we can be reasonably sure there wont be more than 4503599627370495 bits
	end
	local sn=b[1]%52
	for l1=1,math.floor(b[1]/52) do
		table.remove(a,1)
	end
	local o=bn()
	local sns=2^sn
	local snm=2^(53-sn)
	for l1=1,#a do
		o[l1]=math.floor(a[l1]/sns)+(((a[l1+1] or 0)%sns)*snm)
	end
	return o
end

function bn.modexp(b,e,m)
	if m==1 then
		return bn(0)
	end
	local result=bn(1)
	while b>m do
		b=b-m
	end
	while e>0 do
		if e%2==1 then
			result=(result*b)%m
		end
		e=e/2
		b=(b*b)
		while b>m do
			b=b-m
		end
	end
	return result
end

local ffihexnums={[0]=
	0x30,0x31,0x32,0x33,0x34,0x35,0x36,0x37,
	0x38,0x39,0x41,0x42,0x43,0x44,0x45,0x46,
}

local hexnums={
	[0]="0","1","2","3","4","5","6","7",
	"8","9","A","B","C","D","E","F"
}
function bn.tohex(n,noffi)
	n=bn(n,true)
	if ffi and not noffi then
		local o=ffi.new("char[?]",math.ceil((#n)*13)+16)
		local idx
		if n.neg then
			ffi.copy(o,"-0x")
			idx=3
		else
			ffi.copy(o,"0x")
			idx=2
		end
		local st=true
		local c=n[#n] or 0
		for l2=12,0,-1 do
			local h=math.floor(c/(16^l2))%16
			if not st or h~=0 then
				o[idx]=ffihexnums[h]
				idx=idx+1
				st=false
			end
		end
		for l1=#n-1,1,-1 do
			local c=n[l1]
			for l2=12,0,-1 do
				o[idx]=ffihexnums[math.floor(c/(16^l2))%16]
				idx=idx+1
			end
		end
		return ffi.string(o,idx)
	else
		local floor=math.floor
		local o=n.neg and "-0x" or "0x"
		local st=true
		local c=n[#n] or 0
		for l2=12,0,-1 do
			local h=floor(c/(16^l2))%16
			if not st or h~=0 then
				o=o..hexnums[h]
				st=false
			end
		end
		for l1=#n-1,1,-1 do
			local c=n[l1]
			o=o..hexnums[floor(c/281474976710656)%16]..
				hexnums[floor(c/17592186044416)%16]..
				hexnums[floor(c/1099511627776)%16]..
				hexnums[floor(c/68719476736)%16]..
				hexnums[floor(c/4294967296)%16]..
				hexnums[floor(c/268435456)%16]..
				hexnums[floor(c/16777216)%16]..
				hexnums[floor(c/1048576)%16]..
				hexnums[floor(c/65536)%16]..
				hexnums[floor(c/4096)%16]..
				hexnums[floor(c/256)%16]..
				hexnums[floor(c/16)%16]..
				hexnums[c%16]
		end
		return o
	end
end

local function splitfmt(s)
	local o={}
	for c in ("%.0f"):format(s):reverse():gmatch(".") do
		table.insert(o,tonumber(c))
	end
	return o
end

local function carry(n,base)
	local i=1
	while n[i] do
		if n[i]>=base then
			n[i+1]=(n[i+1] or 0)+(math.floor(n[i]/10))
			n[i]=n[i]%base
		end
		i=i+1
	end
	return n
end

function bn.tostring(n)
	n=bn(n,true)
	local q={}
	for l1=1,#n do
		q[l1]=n[l1]
	end
	local c=splitfmt(n[1] or 0)
	local cmul=splitfmt(2^52)
	for l1=2,#n do
		local cn=n[l1]
		for l2=1,52 do
			if (math.floor(cn/(2^(l2-1)))%2)==1 then
				for l1=1,#cmul do
					c[l1]=(c[l1] or 0)+cmul[l1]
				end
				carry(c,10)
			end
			if l1~=#n or l2~=52 then
				for l1=1,#cmul do
					cmul[l1]=cmul[l1]*2
				end
				carry(cmul,10)
			end
		end
	end
	local o=n.neg and "-" or ""
	for l1=1,#c do
		o=c[l1]..o
	end
	return o
end

function bn.tonumber(txt)
	if type(txt)=="number" then
		return bn(txt)
	end
	txt=tostring(txt)
	local o=bn(0)
	local neg
	if txt:sub(1,1)=="-" then
		neg=true
		txt=txt:sub(2)
	end
	local cmul=bn(1)
	for l1=#txt,1,-1 do
		o=o+(cmul*tonumber(txt:sub(l1,l1)))
		if l1~=1 then
			cmul=cmul*10
		end
	end
	o.neg=neg
	return o
end

return bn
