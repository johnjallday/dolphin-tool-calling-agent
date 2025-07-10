-- get_status.lua: print current project status to CLI

function main()
	-- Get current project
	local proj = reaper.EnumProjects(-1, "")

	-- Project path and name
	--local project_path = reaper.GetProjectPath(proj)
	local project_name = reaper.GetProjectName(proj)

	-- Number of tracks
	local track_count = reaper.CountTracks(0)

	-- Output channels
	local num_outputs = reaper.GetNumAudioOutputs()
	local output_names = {}
	for i = 0, num_outputs - 1 do
		local _, name = reaper.GetOutputChannelName(i, "")
		table.insert(output_names, name)
	end

	-- Output latency
	local latency = reaper.GetOutputLatency()

	-- Project tempo (BPM)
	local tempo = reaper.GetSetProjectInfo(0, "PROJECT_TEMPO", 0, false)

	-- Print status to CLI stdout
	--print("Project Path: " .. project_path)
	print("Project Name: " .. project_name)
	print("Track Count: " .. track_count)
	print(string.format("Output Channels (%d):", num_outputs))
	for idx, name in ipairs(output_names) do
		print(string.format("  [%d] %s", idx - 1, name))
	end
	print("Output Latency: " .. latency)
	print("Tempo (BPM): " .. tempo)
end

main()

