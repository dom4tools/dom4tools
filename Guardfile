guard "shell", test: true do
  watch(%r{(.*/)?(.*?)(?:_test)?\.go$}) do |m|
    command = if !m[1].nil? && !m[1].empty?
    "cd #{ m[1] } && go test"
    else
      "go test"
    end

    `#{ command }`
  end
end
