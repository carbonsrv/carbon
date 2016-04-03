-- encoding
encoding = {}
encoding.base64 = {}
function encoding.base64.decode(str)
	if str then
		local data, err = carbon._enc_base64_dec(str)
		if err ~= nil then
			error(err)
		end
		return data
	end
end
function encoding.base64.encode(str)
	if str then
		return carbon._enc_base64_enc(str)
	end
end
