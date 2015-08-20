lock '3.4.0'

set :application, 'minimoozie'
set :repo_url, 'git@github.com:Shopify/minimoozie.git'
set :deploy_to, '/u/apps/minimoozie'
set :keep_releases, 2
set :ejson_file, "config.ejson"
set :ejson_output_file, "config.json"

namespace :docker do
  task :rebuild_image do
    on roles(:app), in: :sequence, wait: 5 do
      execute "cd #{release_path}; docker build -t minimoozie ."
    end
  end
end

namespace :deploy do
  after :restart, :clear_cache do
    execute "docker build -t minimoozie ."
  end
end

after :deploy, "docker:rebuild_image"
