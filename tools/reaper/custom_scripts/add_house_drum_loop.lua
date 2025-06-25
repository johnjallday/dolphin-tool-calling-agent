function add_house_drums_xo()
	-- Determine project root from this script's location
	local info = debug.getinfo(1, "S")
	local scriptPath = info.source:sub(2)
	local scriptDir = scriptPath:match("(.*/)")
	local projectRoot = scriptDir:gsub("/agents/reaper/", "/")
	local templatePath = projectRoot .. "templates/DR_XO_HOUSE_LOOP.RTrackTemplate"
	-- Launch Reaper with the house drums track template
	os.execute('open -a "Reaper" "' .. templatePath .. '"')
end

add_house_drums_xo()
