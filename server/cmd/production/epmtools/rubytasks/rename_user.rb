# frozen_string_literal: true

if ARGV.size == 3
  user = User.find_by(email: ARGV[0])
  if user.nil?
    puts "{ error: \"user not found\" }"
  else
    user.update(first_name: ARGV[1], last_name: ARGV[2])
    puts user.to_json
  end
end
