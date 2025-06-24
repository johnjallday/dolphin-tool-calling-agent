function add_drums_xo_track()
	reaper.InsertTrackAtIndex(reaper.CountTracks(0), true)
	local track = reaper.GetTrack(0, reaper.CountTracks(0) - 1)
	reaper.GetSetMediaTrackInfo_String(track, "P_NAME", "DR Drums XO", true)
	local plugin_name = "AUi: XO (XLN Audio)"
	local ret, insert_fx_idx = reaper.TrackFX_AddByName(track, plugin_name, false, 1)
	if ret < 0 then
		reaper.ShowMessageBox("Failed to insert plugin: " .. plugin_name, "Error", 0)
	end

	reaper.TrackFX_GetOpen(track, 0)
end
-- Function to insert full house drums template
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

-- Decide which function to call based on DRUM_TYPE env var
local drumType = os.getenv("DRUM_TYPE") or ""
if drumType:lower() == "house" then
	add_house_drums_xo()
else
	add_drums_xo_track()
end
