
require 'capistrano_colors'

server "shrimp", :app, :web, :db, primary: true

set (:user) { "USER" }
set (:application) { "shrimp" }
set (:deploy_to) { "/PATH/TO/#{application}" }
set (:nginx_server) { "API.SHRIMP.COM" }
set (:goos) { "linux" }
set (:goarch) { "amd64" }

set :use_sudo, false
set :deploy_via, :remote_cache

ssh_options[:forward_agent] = true

############################

namespace :deploy do

  task :setup do
    config
  end

  task :config do
    run "mkdir -p #{deploy_to}/log"
    run "mkdir -p #{deploy_to}/tmp/pids"
    eye.config
    nginx.config
  end

  task :default do
    build
    update_code
    eye.restart
  end

  desc "Build the program locally"
  task :build, only: { primary: true } do
    run_locally("mkdir -p tmp; GOOS=#{goos} GOARCH=#{goarch} go build -o tmp/#{application}")
  end

  desc "Update the program binary on the host"
  task :update_code, only: [:app] do
    top.upload("tmp/#{application}", "#{deploy_to}/#{application}", {mode: "+x"})
  end

  task :create_symlink do
    logger.info "Nothing to symlink!"
  end

  task :cold do
    build
    update_code
    eye.start
  end
end

############################

namespace :eye do

  %w[start stop restart].each do |a|
    task "#{a}", only: :app do
      run "eye #{a} #{application}"
    end
  end

  task :config do
    template "eye.rb.erb", "/tmp/eye.rb"
    run "mkdir -p /home/#{user}/eye"
    run "mv /tmp/eye.rb /home/#{user}/eye/#{application}.rb"
    run "eye load /home/#{user}/eye/*.rb"
  end
end

############################

namespace :nginx do
  %w[start stop restart].each do |command|
    desc "#{command} nginx"
    task command, roles: :app do
      run "sudo /etc/init.d/nginx #{command}"
    end
  end

  task :config, roles: :app do
    template "nginx.erb", "/tmp/nginx.conf"
    run "#{sudo} mv /tmp/nginx.conf /etc/nginx/sites-enabled/#{application}.conf"
    run "mkdir -p /home/#{user}/www; touch /home/#{user}/www/index.html"
    nginx.restart
  end
end

################################

def template(from, to)
  erb = File.read(File.expand_path("../templates/#{from}", __FILE__))
  put ERB.new(erb).result(binding), to
end