url=$1
while :
do
        u="$(./live -u $url)"
        if [ "$u" = "error" ]
        then
                date
                sleep 5s
                continue
        else
                t="$(./live -t $url)"
        fi
        ffmpeg -i "$u" -c copy "$t.mp4"
done
