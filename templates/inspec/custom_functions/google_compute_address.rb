def address_ip_exists
  !@address.nil?
end

# How many users are there for the address
def user_count
  @users.count
end

# Return the first user resource base name
def user_resource_name
  @users.first.split('/').last
end